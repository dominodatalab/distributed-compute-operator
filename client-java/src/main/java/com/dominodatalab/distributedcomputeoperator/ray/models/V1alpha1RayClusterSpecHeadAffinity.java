/*
 * Kubernetes
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: v1.18.2
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package com.dominodatalab.distributedcomputeoperator.ray.models;

import java.util.Objects;
import java.util.Arrays;
import com.dominodatalab.distributedcomputeoperator.ray.models.V1alpha1RayClusterSpecHeadAffinityNodeAffinity;
import com.dominodatalab.distributedcomputeoperator.ray.models.V1alpha1RayClusterSpecHeadAffinityPodAffinity;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;

/**
 * V1alpha1RayClusterSpecHeadAffinity
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2021-10-04T16:37:47.111Z[Etc/UTC]")
public class V1alpha1RayClusterSpecHeadAffinity {
  public static final String SERIALIZED_NAME_NODE_AFFINITY = "nodeAffinity";
  @SerializedName(SERIALIZED_NAME_NODE_AFFINITY)
  private V1alpha1RayClusterSpecHeadAffinityNodeAffinity nodeAffinity;

  public static final String SERIALIZED_NAME_POD_AFFINITY = "podAffinity";
  @SerializedName(SERIALIZED_NAME_POD_AFFINITY)
  private V1alpha1RayClusterSpecHeadAffinityPodAffinity podAffinity;

  public static final String SERIALIZED_NAME_POD_ANTI_AFFINITY = "podAntiAffinity";
  @SerializedName(SERIALIZED_NAME_POD_ANTI_AFFINITY)
  private V1alpha1RayClusterSpecHeadAffinityPodAffinity podAntiAffinity;


  public V1alpha1RayClusterSpecHeadAffinity nodeAffinity(V1alpha1RayClusterSpecHeadAffinityNodeAffinity nodeAffinity) {
    
    this.nodeAffinity = nodeAffinity;
    return this;
  }

   /**
   * Get nodeAffinity
   * @return nodeAffinity
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public V1alpha1RayClusterSpecHeadAffinityNodeAffinity getNodeAffinity() {
    return nodeAffinity;
  }


  public void setNodeAffinity(V1alpha1RayClusterSpecHeadAffinityNodeAffinity nodeAffinity) {
    this.nodeAffinity = nodeAffinity;
  }


  public V1alpha1RayClusterSpecHeadAffinity podAffinity(V1alpha1RayClusterSpecHeadAffinityPodAffinity podAffinity) {
    
    this.podAffinity = podAffinity;
    return this;
  }

   /**
   * Get podAffinity
   * @return podAffinity
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public V1alpha1RayClusterSpecHeadAffinityPodAffinity getPodAffinity() {
    return podAffinity;
  }


  public void setPodAffinity(V1alpha1RayClusterSpecHeadAffinityPodAffinity podAffinity) {
    this.podAffinity = podAffinity;
  }


  public V1alpha1RayClusterSpecHeadAffinity podAntiAffinity(V1alpha1RayClusterSpecHeadAffinityPodAffinity podAntiAffinity) {
    
    this.podAntiAffinity = podAntiAffinity;
    return this;
  }

   /**
   * Get podAntiAffinity
   * @return podAntiAffinity
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public V1alpha1RayClusterSpecHeadAffinityPodAffinity getPodAntiAffinity() {
    return podAntiAffinity;
  }


  public void setPodAntiAffinity(V1alpha1RayClusterSpecHeadAffinityPodAffinity podAntiAffinity) {
    this.podAntiAffinity = podAntiAffinity;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    V1alpha1RayClusterSpecHeadAffinity v1alpha1RayClusterSpecHeadAffinity = (V1alpha1RayClusterSpecHeadAffinity) o;
    return Objects.equals(this.nodeAffinity, v1alpha1RayClusterSpecHeadAffinity.nodeAffinity) &&
        Objects.equals(this.podAffinity, v1alpha1RayClusterSpecHeadAffinity.podAffinity) &&
        Objects.equals(this.podAntiAffinity, v1alpha1RayClusterSpecHeadAffinity.podAntiAffinity);
  }

  @Override
  public int hashCode() {
    return Objects.hash(nodeAffinity, podAffinity, podAntiAffinity);
  }


  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class V1alpha1RayClusterSpecHeadAffinity {\n");
    sb.append("    nodeAffinity: ").append(toIndentedString(nodeAffinity)).append("\n");
    sb.append("    podAffinity: ").append(toIndentedString(podAffinity)).append("\n");
    sb.append("    podAntiAffinity: ").append(toIndentedString(podAntiAffinity)).append("\n");
    sb.append("}");
    return sb.toString();
  }

  /**
   * Convert the given object to string with each line indented by 4 spaces
   * (except the first line).
   */
  private String toIndentedString(Object o) {
    if (o == null) {
      return "null";
    }
    return o.toString().replace("\n", "\n    ");
  }

}

